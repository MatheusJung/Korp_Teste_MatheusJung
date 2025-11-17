import { Component ,Input} from '@angular/core';

@Component({
  selector: 'app-button',
  imports: [],
  templateUrl: './button.html',
  styleUrl: './button.scss',
})
export class Button {
  @Input() title:string = 'Título indisponível';
  @Input() imageUrl:string = 'image/default_photo.svg';
  @Input() linkUrl:string = '';
}
